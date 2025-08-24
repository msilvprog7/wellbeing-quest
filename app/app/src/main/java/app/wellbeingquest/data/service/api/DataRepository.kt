package app.wellbeingquest.data.service.api

import android.util.Log
import app.wellbeingquest.data.local.database.AppDatabase
import app.wellbeingquest.data.local.entity.EntryQueueItem
import app.wellbeingquest.data.local.entity.SuggestionCacheItem
import app.wellbeingquest.data.local.entity.WeekCacheItem
import app.wellbeingquest.data.model.Activity
import app.wellbeingquest.data.model.Feeling
import app.wellbeingquest.data.model.Suggestions
import app.wellbeingquest.data.service.dto.ActivityDto
import java.time.Instant
import java.time.LocalDate
import java.time.OffsetDateTime
import java.time.ZoneId
import java.time.ZoneOffset
import java.time.format.DateTimeFormatter
import java.time.format.DateTimeParseException

// todo: implement week and suggestion data functions between api service and app database
class DataRepository(
    private val appDatabase: AppDatabase,
    private val apiService: ApiService
) {
    suspend fun getEntriesToUpload(): List<EntryQueueItem> {
        return appDatabase.entryQueueItemDao().getQueueItems()
    }

    suspend fun uploadEntry(entry: EntryQueueItem) {
        // call api service to upload entry
        val response = apiService.postActivity(
            ActivityDto(
                name = entry.activity,
                feelings = entry.feelings
            )
        )

        if (!response.isSuccessful || response.body() == null) {
            throw Exception("Failed to upload entry, response code: ${response.code()}, message: ${response.message()}")
        }

        // add entry to database
        val activity = response.body()!!
        if (activity.week == null ||
            activity.created == null ||
            activity.feelings == null
        ) {
            throw Exception("Failed to upload entry, activity in response invalid, response code: ${response.code()}, message: ${response.message()}")
        }

        for (feeling in activity.feelings!!) {
            appDatabase.weekCacheItemDao().insert(
                WeekCacheItem(
                    name = activity.week!!,
                    feeling = feeling,
                    activity = activity.name,
                    created = parseLocalDateToMs(activity.created!!)
                )
            )
        }

        // remove entry from local database
        appDatabase.entryQueueItemDao().delete(entry.id)
    }

    suspend fun getWeek(week: String): Map<String, Feeling> {
        val feelings = HashMap<String, Feeling>()

        // Sync from api
        try {
            syncWeek(week)
        } catch (e: Exception) {
            Log.e("DataRepository", "Failed to sync week from API", e)
        }

        // Add from cache
        val weekCacheItems = appDatabase.weekCacheItemDao().getItemsByWeek(week)
        for (item in weekCacheItems) {
            val feeling = feelings.getOrPut(item.feeling, { Feeling(item.feeling, mutableListOf()) })
            feeling.activities.add(Activity(item.activity, true, true))
        }

        // Add from upload queue
        val entryQueueItems = appDatabase.entryQueueItemDao().getQueueItems()
        for (item in entryQueueItems) {
            val itemWeek = getWeek(getStartOfWeekFromLocalDate(getLocalDate(item.created)))
            if (week != itemWeek) {
                continue;
            }

            for (feeling in item.feelings) {
                val feeling = feelings.getOrPut(feeling, { Feeling(feeling, mutableListOf()) })
                feeling.activities.add(Activity(item.activity, false, true))
            }
        }

        return feelings
    }

    suspend fun syncWeek(week: String) {
        val response = apiService.getWeek(week)

        if (!response.isSuccessful && response.code() != 404) {
            throw Exception("Failed to sync week from API, response code: ${response.code()}, message: ${response.message()}")
        }

        // Clear cache for week
        var dao = appDatabase.weekCacheItemDao()
        dao.deleteByWeek(week)

        if (response.code() == 404 || response.body() == null || response.body()!!.feelings == null) {
            return;
        }

        // Update cache
        for (feeling in response.body()!!.feelings) {
            if (feeling.activities == null) {
                continue;
            }

            for (activity in feeling.activities) {
                dao.insert(
                    WeekCacheItem(
                        name = week,
                        feeling = feeling.name,
                        activity = activity.name,
                        created = parseLocalDateToMs(activity.created)
                    )
                )
            }
        }
    }

    suspend fun getSuggestions(): Suggestions {
        val activities = mutableListOf<String>()
        val feelings = mutableListOf<String>()

        // Sync from api
        try {
            syncSuggestions()
        } catch (e: Exception) {
            Log.e("DataRepository", "Failed to sync suggestions from API", e)
        }

        // Add from cache
        var dao = appDatabase.suggestionCacheItemDao()
        val items = dao.getItems()
        for (item in items) {
            when (item.type) {
                "activity" -> {
                    activities.add(item.text)
                }
                "feeling" -> {
                    feelings.add(item.text)
                }
                else -> {
                    Log.e("DataRepository", "Unknown suggestion type: ${item.type}")
                }
            }
        }

        return Suggestions(activities.toList(), feelings.toList())
    }

    suspend fun syncSuggestions() {
        val response = apiService.getSuggestions()

        if (!response.isSuccessful && response.code() != 404) {
            throw Exception("Failed to sync suggestions from API, response code: ${response.code()}, message: ${response.message()}")
        }

        // Clear cache for week
        var dao = appDatabase.suggestionCacheItemDao()
        dao.deleteAll()

        if (response.code() == 404 ||
            response.body() == null ||
            response.body()!!.activities == null ||
            response.body()!!.feelings == null) {
            return;
        }

        // Update cache
        for (activity in response.body()!!.activities) {
            dao.insert(SuggestionCacheItem(
                type = "activity",
                text = activity.name
            ))
        }

        for (feeling in response.body()!!.feelings) {
            dao.insert(SuggestionCacheItem(
                type = "feeling",
                text = feeling.name
            ))
        }
    }

    private fun parseLocalDateToMs(dateTimeString: String?): Long {
        if (dateTimeString.isNullOrBlank()) {
            throw IllegalArgumentException("Date time string cannot be null or blank.")
        }

        return try {
            val formatter = DateTimeFormatter.ISO_OFFSET_DATE_TIME
            val offsetDateTime = OffsetDateTime.parse(dateTimeString, formatter)
            val localDate = offsetDateTime.toLocalDate()
            val startOfDayUtc = localDate.atStartOfDay(ZoneOffset.UTC)
            startOfDayUtc.toInstant().toEpochMilli()
        } catch (e: DateTimeParseException) {
            // Re-throw with a more informative message or handle as needed
            throw DateTimeParseException(
                "Failed to parse date string: '$dateTimeString'. Expected ISO_OFFSET_DATE_TIME format (e.g., 'yyyy-MM-ddTHH:mm:ss.SSSSSSSSSXXX').",
                e.parsedString,
                e.errorIndex,
                e
            )
        }
    }

    private fun getStartOfWeekFromLocalDate(date: LocalDate): LocalDate {
        val daysAfterSunday = date.dayOfWeek.value % 7
        return date.minusDays(daysAfterSunday.toLong())
    }

    private fun getWeek(start: LocalDate): String {
        return DateTimeFormatter.ISO_LOCAL_DATE.format(start)
    }

    private fun getLocalDate(ms: Long): LocalDate {
        return Instant.ofEpochMilli(ms).atZone(ZoneId.systemDefault()).toLocalDate()
    }
}