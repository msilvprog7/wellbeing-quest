package app.wellbeingquest.data.local.dao

import androidx.room.Dao
import androidx.room.Insert
import androidx.room.OnConflictStrategy
import androidx.room.Query
import app.wellbeingquest.data.local.entity.EntryEntity

@Dao
interface EntryDao {
    @Query("SELECT * FROM entries WHERE weekId = :weekId")
    fun findByWeekId(weekId: Int): List<EntryEntity>

    @Insert(onConflict = OnConflictStrategy.ABORT)
    suspend fun insert(entry: EntryEntity): Int
}