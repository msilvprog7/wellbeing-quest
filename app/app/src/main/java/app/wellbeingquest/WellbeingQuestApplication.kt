package app.wellbeingquest

import android.app.Application

class WellbeingQuestApplication : Application() {
    companion object {
        lateinit var instance: WellbeingQuestApplication
            private set
    }

    override fun onCreate() {
        super.onCreate()
        instance = this
    }
}