package app.wellbeingquest.data.service.api

import android.content.Context
import android.util.Log
import androidx.hilt.work.HiltWorker
import androidx.work.BackoffPolicy
import androidx.work.Constraints
import androidx.work.CoroutineWorker
import androidx.work.Data
import androidx.work.ExistingWorkPolicy
import androidx.work.NetworkType
import androidx.work.OneTimeWorkRequest
import androidx.work.OneTimeWorkRequestBuilder
import androidx.work.WorkManager
import androidx.work.WorkerParameters
import app.wellbeingquest.data.local.database.DatabaseProvider
import dagger.assisted.Assisted
import dagger.assisted.AssistedInject
import java.util.concurrent.TimeUnit
import javax.inject.Inject

fun scheduleUploadWorker(context: Context): OneTimeWorkRequest {
    Log.d("UploadWorker", "Scheduling upload worker")

    val constraints = Constraints.Builder()
        .setRequiredNetworkType(NetworkType.CONNECTED)
        .build()

    val uploadWorkRequest = OneTimeWorkRequestBuilder<UploadWorker>()
        .setConstraints(constraints)
        .setBackoffCriteria(
            BackoffPolicy.EXPONENTIAL,
            10,
            TimeUnit.SECONDS
        )
        .build()

    WorkManager.getInstance(context).enqueueUniqueWork(
        UploadWorker.WORK_NAME,
        ExistingWorkPolicy.REPLACE,
        uploadWorkRequest
    )

    return uploadWorkRequest
}

@HiltWorker
class UploadWorker @AssistedInject constructor(
    @Assisted appContext: Context,
    @Assisted workerParams: WorkerParameters,
    private val api: ApiService,
    private val dataRepository: DataRepository
) : CoroutineWorker(appContext, workerParams) {

    override suspend fun doWork(): Result {
        Log.d("UploadWorker", "Uploading entries")

        try {
            val entries = dataRepository.getEntriesToUpload()
            Log.d("UploadWorker", "Entries to upload: ${entries.size}")

            for (entry in entries) {
                dataRepository.uploadEntry(entry)
            }

            return Result.success()
        } catch (e: Exception) {
            Log.e("UploadWorker", "Error uploading entries", e)
            return Result.failure()
        }
    }

    companion object {
        const val WORK_NAME = "UploadWorker"
    }

}