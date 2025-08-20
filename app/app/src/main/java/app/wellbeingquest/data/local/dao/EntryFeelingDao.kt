package app.wellbeingquest.data.local.dao

import androidx.room.Dao
import androidx.room.Insert
import androidx.room.OnConflictStrategy
import androidx.room.Query
import app.wellbeingquest.data.local.entity.EntryFeelingEntity

@Dao
interface EntryFeelingDao {
    @Query("SELECT * FROM entryFeelings WHERE entryId = :entryId")
    fun findByEntryId(entryId: Int): List<EntryFeelingEntity>

    @Query("SELECT * FROM entryFeelings WHERE feelingId = :feelingId")
    fun findByFeelingId(feelingId: Int): List<EntryFeelingEntity>

    @Query("SELECT * FROM entryFeelings WHERE entryId = :activityId AND feelingId = :feelingId")
    fun findByIds(activityId: Int, feelingId: Int): EntryFeelingEntity?

    @Insert(onConflict = OnConflictStrategy.ABORT)
    suspend fun insert(entryFeeling: EntryFeelingEntity)
}