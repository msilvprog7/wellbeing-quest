package app.wellbeingquest.data.service.api

import app.wellbeingquest.data.local.database.AppDatabase
import app.wellbeingquest.data.local.entity.EntryQueueItem

// todo: upload entry to api service
// todo: remove entry from local database
// todo: implement week and suggestion data functions between api service and app database
class DataRepository(
    private val appDatabase: AppDatabase,
    private val apiService: ApiService
){
    suspend fun getEntriesToUpload(): List<EntryQueueItem> {
        return appDatabase.entryQueueItemDao().getQueueItems()
    }

    suspend fun uploadEntry(entry: EntryQueueItem) {
    }
}