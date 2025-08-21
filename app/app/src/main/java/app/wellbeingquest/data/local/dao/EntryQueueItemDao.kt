package app.wellbeingquest.data.local.dao

import androidx.room.Dao
import androidx.room.Insert
import androidx.room.OnConflictStrategy
import androidx.room.Query
import app.wellbeingquest.data.local.entity.EntryQueueItem

@Dao
interface EntryQueueItemDao {
    @Query("SELECT * FROM entryQueueItems ORDER BY created ASC")
    suspend fun getQueueItems(): List<EntryQueueItem>

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insert(entry: EntryQueueItem)

    @Query("DELETE FROM entryQueueItems WHERE id = :queueItemId")
    suspend fun delete(queueItemId: Int)
}
