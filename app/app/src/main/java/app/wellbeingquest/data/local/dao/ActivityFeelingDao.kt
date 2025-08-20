package app.wellbeingquest.data.local.dao

import androidx.room.Dao
import androidx.room.Insert
import androidx.room.OnConflictStrategy
import androidx.room.Query
import app.wellbeingquest.data.local.entity.ActivityFeelingEntity

@Dao
interface ActivityFeelingDao {
    @Query("SELECT * FROM activityFeelings WHERE activityId = :activityId")
    fun findByActivityId(activityId: Int): List<ActivityFeelingEntity>

    @Query("SELECT * FROM activityFeelings WHERE feelingId = :feelingId")
    fun findByFeelingId(feelingId: Int): List<ActivityFeelingEntity>

    @Query("SELECT * FROM activityFeelings WHERE activityId = :activityId AND feelingId = :feelingId")
    fun findByIds(activityId: Int, feelingId: Int): ActivityFeelingEntity?

    @Insert(onConflict = OnConflictStrategy.ABORT)
    suspend fun insert(activityFeeling: ActivityFeelingEntity)
}