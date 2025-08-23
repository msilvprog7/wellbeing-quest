package app.wellbeingquest.data.service.api

import app.wellbeingquest.data.local.database.AppDatabase
import app.wellbeingquest.data.local.entity.EntryQueueItem
import app.wellbeingquest.data.local.entity.WeekCacheItem
import app.wellbeingquest.data.service.dto.ActivityDto
import java.time.LocalDate
import java.time.OffsetDateTime
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
}