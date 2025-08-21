package app.wellbeingquest.data.local.dao

import androidx.room.Dao
import androidx.room.Insert
import androidx.room.OnConflictStrategy
import androidx.room.Query
import app.wellbeingquest.data.local.entity.EntryDraft

@Dao
interface EntryDraftDao {
    @Query("SELECT * FROM entryDrafts LIMIT 1")
    suspend fun getDraft(): EntryDraft?

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insert(entry: EntryDraft)

    @Query("DELETE FROM entryDrafts")
    suspend fun delete()
}
