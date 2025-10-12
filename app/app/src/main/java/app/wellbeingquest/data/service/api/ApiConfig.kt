package app.wellbeingquest.data.service.api

import android.content.Context
import app.wellbeingquest.BuildConfig
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.runBlocking

object ApiConfig {

    private lateinit var contextRef: Context

    var sessionToken: String
        get() = runBlocking {
            DataStoreManager.getSessionToken(contextRef).first()
        }
        set(value) {
            runBlocking {
                DataStoreManager.setSessionToken(contextRef, value)
            }
        }

    fun init(context: Context) {
        contextRef = context.applicationContext
    }

    fun getBaseUrl(context: Context): String {
        val prefs = context.getSharedPreferences("app_prefs", Context.MODE_PRIVATE)
        return prefs.getString("base_url", BuildConfig.SERVICE_URL)!!
    }
}
