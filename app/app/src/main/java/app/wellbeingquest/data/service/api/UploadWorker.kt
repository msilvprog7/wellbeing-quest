package app.wellbeingquest.data.service.api

import android.content.Context
import androidx.work.CoroutineWorker
import androidx.work.WorkerParameters

// todo: upload EntryQueueItems
// todo: trigger worker when items added and app started
class UploadWorker(
    appContext: Context,
    workerParams: WorkerParameters,
    private val dataRepository: DataRepository
) : CoroutineWorker(appContext, workerParams) {
    override suspend fun doWork(): Result {
        return Result.success()
    }
}