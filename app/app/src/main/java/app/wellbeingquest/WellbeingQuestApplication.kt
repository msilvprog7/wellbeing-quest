package app.wellbeingquest

import android.app.Application
import androidx.hilt.work.HiltWorkerFactory
import androidx.work.Configuration
import dagger.hilt.android.HiltAndroidApp
import javax.inject.Inject

@HiltAndroidApp
class WellbeingQuestApplication : Application(), Configuration.Provider {

    @Inject
    lateinit var workerFactory: HiltWorkerFactory

    companion object {
        lateinit var instance: WellbeingQuestApplication
            private set
    }

    override val workManagerConfiguration: Configuration
            get() = Configuration.Builder()
                .setWorkerFactory(workerFactory)
                .setMinimumLoggingLevel(loggingLevel = android.util.Log.DEBUG)
                .build()
}