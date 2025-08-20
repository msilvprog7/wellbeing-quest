package app.wellbeingquest.data.local.dao

import androidx.room.Dao
import androidx.room.Insert
import androidx.room.OnConflictStrategy
import androidx.room.Query
import app.wellbeingquest.data.local.entity.FeelingEntity

@Dao
interface FeelingDao {
    @Query("SELECT * FROM feelings WHERE name = :name LIMIT 1")
    fun findByName(name: String): FeelingEntity?

    @Insert(onConflict = OnConflictStrategy.ABORT)
    suspend fun insert(feeling: FeelingEntity): Int
}