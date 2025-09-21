package app.wellbeingquest.data.service.api

import android.content.Context

object ApiConfig {
    fun getBaseUrl(context: Context): String {
        val prefs = context.getSharedPreferences("app_prefs", Context.MODE_PRIVATE)
        return prefs.getString("base_url", "https://api.wellbeingquest.app/")!!
    }
}
