package app.wellbeingquest

import android.app.Application
import app.wellbeingquest.data.service.api.ApiConfig

class WellbeingQuestApplication : Application() {
    companion object {
        lateinit var instance: WellbeingQuestApplication
            private set
    }

    override fun onCreate() {
        super.onCreate()
        ApiConfig.init(this)
        instance = this
    }
}