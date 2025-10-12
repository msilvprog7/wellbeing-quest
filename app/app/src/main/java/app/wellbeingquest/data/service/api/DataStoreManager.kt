package app.wellbeingquest.data.service.api

import android.content.Context
import androidx.datastore.preferences.core.edit
import androidx.datastore.preferences.core.stringPreferencesKey
import androidx.datastore.preferences.preferencesDataStore
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.map

private const val DATASTORE_NAME = "secure_prefs"
val Context.dataStore by preferencesDataStore(name = DATASTORE_NAME)

object DataStoreManager {
    private val SESSION_TOKEN = stringPreferencesKey("session_token")

    fun getSessionToken(context: Context): Flow<String> {
        val crypto = CryptoManager(context)
        return context.dataStore.data.map { prefs ->
            prefs[SESSION_TOKEN]?.let { crypto.decrypt(it) } ?: ""
        }
    }

    suspend fun setSessionToken(context: Context, token: String) {
        val crypto = CryptoManager(context)
        val encrypted = crypto.encrypt(token)
        context.dataStore.edit { it[SESSION_TOKEN] = encrypted }
    }
}