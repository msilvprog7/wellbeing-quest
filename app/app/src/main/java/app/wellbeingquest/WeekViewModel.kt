package app.wellbeingquest

import android.content.Context
import android.util.Log
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import androidx.work.WorkInfo
import androidx.work.WorkManager
import app.wellbeingquest.data.model.Feeling
import app.wellbeingquest.data.service.api.DataRepository
import app.wellbeingquest.data.service.api.scheduleUploadWorker
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.SharingStarted
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.map
import kotlinx.coroutines.flow.stateIn
import kotlinx.coroutines.launch
import java.time.LocalDate
import java.time.format.DateTimeFormatter

class WeekViewModel(
    private val context: Context,
    private val dataRepository: DataRepository
): ViewModel() {
    private val _currentWeekActualStart = getStartOfWeek()
    private val _selectedWeekStart = MutableStateFlow(_currentWeekActualStart)
    private val _feelingsInWeek = MutableStateFlow(listOf<Feeling>())
    val selectedWeekStart: StateFlow<LocalDate> = _selectedWeekStart
    val hasNextWeek: StateFlow<Boolean> = _selectedWeekStart
        .map { selectedDate ->
            selectedDate.isBefore(_currentWeekActualStart)
        }
        .stateIn(viewModelScope, SharingStarted.Eagerly, _selectedWeekStart.value.isBefore(_currentWeekActualStart))
    val feelingsInWeek: StateFlow<List<Feeling>> = _feelingsInWeek

    init {
        triggerUploadWorker()
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

    fun triggerUploadWorker() {
        val request = scheduleUploadWorker(context)
        viewModelScope.launch {
            WorkManager.getInstance(context)
                .getWorkInfoByIdFlow(request.id)
                .collect { workInfo ->
                    Log.d("WeekViewModel", "UploadWorker State: ${workInfo?.state}")

                    if (workInfo?.state == WorkInfo.State.SUCCEEDED) {
                        Log.d("WeekViewModel", "UploadWorker SUCCEEDED. Refreshing week data.")
                        loadWeek(_selectedWeekStart.value)
                    } else {
                        Log.e("WeekViewModel", "UploadWorker FAILED.")
                    }
                }
        }
    }

    private fun getStartOfWeek(): LocalDate {
        return getStartOfWeekFromLocalDate(LocalDate.now())
    }

    private fun getStartOfWeekFromLocalDate(date: LocalDate): LocalDate {
        val daysAfterSunday = date.dayOfWeek.value % 7
        return date.minusDays(daysAfterSunday.toLong())
    }

    private fun getWeek(start: LocalDate): String {
        return DateTimeFormatter.ISO_LOCAL_DATE.format(start)
    }

    private fun loadWeek(start: LocalDate) {
        val week = getWeek(start)

        viewModelScope.launch {
            val feelings = dataRepository.getWeek(week)
            _feelingsInWeek.value = feelings.values.toList()
        }
    }
}