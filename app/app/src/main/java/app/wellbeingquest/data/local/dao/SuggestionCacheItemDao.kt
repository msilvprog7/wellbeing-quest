package app.wellbeingquest.data.local.dao

import androidx.room.Dao
import androidx.room.Insert
import androidx.room.OnConflictStrategy
import androidx.room.Query
import app.wellbeingquest.data.local.entity.SuggestionCacheItem

@Dao
interface SuggestionCacheItemDao {
    @Query("SELECT * FROM suggestionCacheItems ORDER BY created ASC")
    suspend fun getItems(): List<SuggestionCacheItem>

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insert(item: SuggestionCacheItem)

    @Query("DELETE FROM suggestionCacheItems")
    suspend fun deleteAll()
}
