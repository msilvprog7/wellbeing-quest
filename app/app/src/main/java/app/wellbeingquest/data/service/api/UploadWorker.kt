package app.wellbeingquest.data.service.api

import android.content.Context
import android.util.Log
import androidx.work.BackoffPolicy
import androidx.work.Constraints
import androidx.work.CoroutineWorker
import androidx.work.ExistingWorkPolicy
import androidx.work.NetworkType
import androidx.work.OneTimeWorkRequest
import androidx.work.OneTimeWorkRequestBuilder
import androidx.work.WorkManager
import androidx.work.WorkerParameters
import app.wellbeingquest.data.local.database.DatabaseProvider
import retrofit2.Retrofit
import java.util.concurrent.TimeUnit

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

class UploadWorker(
    appContext: Context,
    workerParams: WorkerParameters
) : CoroutineWorker(appContext, workerParams) {

    override suspend fun doWork(): Result {
        Log.d("UploadWorker", "Uploading entries")

        try {
            // todo: setup hilt to inject dependencies throughout app
            val appDatabase = DatabaseProvider.getInstance(applicationContext)
            val apiService = RetrofitInstance.api
            val dataRepository = DataRepository(appDatabase, apiService)

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