package app.wellbeingquest

import android.content.Context
import android.content.SharedPreferences
import android.util.Log
import androidx.datastore.core.DataStore
import androidx.datastore.preferences.core.Preferences
import androidx.datastore.preferences.core.edit
import androidx.datastore.preferences.core.stringPreferencesKey
import app.wellbeingquest.data.service.api.CryptoManager
import dagger.hilt.android.qualifiers.ApplicationContext
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.map
import kotlinx.coroutines.launch
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class SessionManager @Inject constructor(
    @ApplicationContext private val context: Context,
    private val dataStore: DataStore<Preferences>,
    private val cryptoManager: CryptoManager,
    private val applicationScope: CoroutineScope
) {

    @Volatile
    private var cachedToken: String? = null

    companion object {
        private val SESSION_TOKEN_KEY = stringPreferencesKey("session_token")
        private const val APP_PREFS_NAME = "app_prefs"
    }

    init {
        applicationScope.launch {
            sessionTokenFlow.collect { token ->
                cachedToken = token
            }
        }
    }

    val sessionTokenFlow: Flow<String?> = dataStore.data.map { preferences ->
        try {
            val encryptedToken = preferences[SESSION_TOKEN_KEY]

            if (encryptedToken == null) {
                return@map null
            }

            cryptoManager.decrypt(encryptedToken)
        } catch (e: javax.crypto.AEADBadTagException) {
            Log.e("SessionManager", "Failed to decrypt token (AEADBadTagException). Key likely invalid/lost.", e)
            null
        } catch (e: Exception) {
            Log.e("SessionManager", "An unexpected error occurred during token decryption.", e)
            null
        }
    }

    fun getCachedSessionToken(): String? = cachedToken

    suspend fun saveSessionToken(token: String) {
        dataStore.edit { preferences ->
            val encryptedToken = cryptoManager.encrypt(token)
            preferences[SESSION_TOKEN_KEY] = encryptedToken
        }
    }

    suspend fun clearSessionToken() {
        dataStore.edit { preferences ->
            preferences.remove(SESSION_TOKEN_KEY)
        }
    }

    fun getBaseUrl(): String {
        val prefs: SharedPreferences = context.getSharedPreferences(APP_PREFS_NAME, Context.MODE_PRIVATE)
        return prefs.getString("base_url", BuildConfig.SERVICE_URL) ?: BuildConfig.SERVICE_URL
    }
}