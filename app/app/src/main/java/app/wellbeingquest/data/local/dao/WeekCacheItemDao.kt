package app.wellbeingquest.data.local.dao

import androidx.room.Dao
import androidx.room.Insert
import androidx.room.OnConflictStrategy
import androidx.room.Query
import app.wellbeingquest.data.local.entity.WeekCacheItem

@Dao
interface WeekCacheItemDao {
    @Query("SELECT * FROM weekCacheItems WHERE name = :name ORDER BY created ASC")
    suspend fun getItemsByWeek(name: String): List<WeekCacheItem>

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insert(item: WeekCacheItem)
}
