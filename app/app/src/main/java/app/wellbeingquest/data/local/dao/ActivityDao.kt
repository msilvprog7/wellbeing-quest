package app.wellbeingquest.data.local.dao

import androidx.room.Dao
import androidx.room.Insert
import androidx.room.OnConflictStrategy
import androidx.room.Query
import app.wellbeingquest.data.local.entity.ActivityEntity

@Dao
interface ActivityDao {
    @Query("SELECT * FROM activities WHERE name = :name LIMIT 1")
    fun findByName(name: String): ActivityEntity?

    @Insert(onConflict = OnConflictStrategy.ABORT)
    suspend fun insert(activity: ActivityEntity): Int
}