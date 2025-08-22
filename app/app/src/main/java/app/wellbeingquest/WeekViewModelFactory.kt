package app.wellbeingquest

import androidx.lifecycle.ViewModel
import androidx.lifecycle.ViewModelProvider
import app.wellbeingquest.data.local.database.AppDatabase

class WeekViewModelFactory(private val appDatabase: AppDatabase) : ViewModelProvider.Factory {
    override fun <T : ViewModel> create(modelClass: Class<T>): T {
        if (modelClass.isAssignableFrom(WeekViewModel::class.java)) {
            @Suppress("UNCHECKED_CAST")
            return WeekViewModel(appDatabase) as T
        }

        throw IllegalArgumentException("Unknown ViewModel class")
    }
}
