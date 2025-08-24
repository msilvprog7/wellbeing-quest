package app.wellbeingquest

import android.content.Context
import androidx.lifecycle.ViewModel
import androidx.lifecycle.ViewModelProvider
import app.wellbeingquest.data.service.api.DataRepository

class WeekViewModelFactory(
    private val context: Context,
    private val dataRepository: DataRepository
) : ViewModelProvider.Factory {
    override fun <T : ViewModel> create(modelClass: Class<T>): T {
        if (modelClass.isAssignableFrom(WeekViewModel::class.java)) {
            @Suppress("UNCHECKED_CAST")
            return WeekViewModel(context, dataRepository) as T
        }

        throw IllegalArgumentException("Unknown ViewModel class")
    }
}
