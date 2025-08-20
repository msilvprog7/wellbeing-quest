package app.wellbeingquest.data.local.dao

import androidx.room.Dao
import androidx.room.Insert
import androidx.room.OnConflictStrategy
import androidx.room.Query
import app.wellbeingquest.data.local.entity.WeekEntity

@Dao
interface WeekDao {
    @Query("SELECT * FROM weeks WHERE name = :name LIMIT 1")
    fun findByName(name: String): WeekEntity?

    @Insert(onConflict = OnConflictStrategy.ABORT)
    suspend fun insert(week: WeekEntity): Int
}