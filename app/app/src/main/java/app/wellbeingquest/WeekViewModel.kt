package app.wellbeingquest

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.SharingStarted
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.map
import kotlinx.coroutines.flow.stateIn
import java.time.LocalDate

class WeekViewModel : ViewModel() {

    private val _selectedWeekStart = MutableStateFlow(getStartOfWeek())
    val selectedWeekStart: StateFlow<LocalDate> = _selectedWeekStart
    val hasNextWeek: StateFlow<Boolean> = _selectedWeekStart
        .map { date -> date.plusWeeks(1).isBefore(LocalDate.now()) }
        .stateIn(viewModelScope, SharingStarted.Eagerly, false)

    fun previousWeek() {
        _selectedWeekStart.value = _selectedWeekStart.value.minusWeeks(1)
    }

    fun nextWeek() {
        if (!hasNextWeek.value) {
            return
        }

        _selectedWeekStart.value = _selectedWeekStart.value.plusWeeks(1)
    }

    fun getStartOfWeek(): LocalDate {
        val today = LocalDate.now()
        val daysAfterSunday = today.dayOfWeek.value % 7
        return today.minusDays(daysAfterSunday.toLong())
    }
}