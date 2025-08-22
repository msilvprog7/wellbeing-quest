package app.wellbeingquest

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import app.wellbeingquest.data.local.database.AppDatabase
import app.wellbeingquest.data.model.Activity
import app.wellbeingquest.data.model.Feeling
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.SharingStarted
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.map
import kotlinx.coroutines.flow.stateIn
import kotlinx.coroutines.launch
import java.time.Instant
import java.time.LocalDate
import java.time.ZoneId
import java.time.format.DateTimeFormatter

class WeekViewModel(
    private val appDatabase: AppDatabase
): ViewModel() {
    private val _currentWeekStart = getStartOfWeek()
    private val _currentWeek = getWeek(_currentWeekStart)
    private val _selectedWeekStart = MutableStateFlow(_currentWeekStart)
    private val _feelingsInWeek = MutableStateFlow(listOf<Feeling>())
    val selectedWeekStart: StateFlow<LocalDate> = _selectedWeekStart
    val hasNextWeek: StateFlow<Boolean> = _selectedWeekStart
        .map { date -> date.plusWeeks(1).isBefore(LocalDate.now()) }
        .stateIn(viewModelScope, SharingStarted.Eagerly, false)
    val feelingsInWeek: StateFlow<List<Feeling>> = _feelingsInWeek

    init {
        loadWeek(_selectedWeekStart.value)
    }

    fun previousWeek() {
        _selectedWeekStart.value = _selectedWeekStart.value.minusWeeks(1)
        loadWeek(_selectedWeekStart.value)
    }

    fun nextWeek() {
        if (!hasNextWeek.value) {
            return
        }

        _selectedWeekStart.value = _selectedWeekStart.value.plusWeeks(1)
        loadWeek(_selectedWeekStart.value)
    }

    fun getStartOfWeek(): LocalDate {
        return getStartOfWeekFromLocalDate(LocalDate.now())
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

    private fun loadWeek(start: LocalDate) {
        val week = getWeek(start)

        viewModelScope.launch {
            val feelings = HashMap<String, Feeling>()

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

            _feelingsInWeek.value = feelings.values.toList()
        }
    }
}