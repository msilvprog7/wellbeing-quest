package app.wellbeingquest.data.local.entity

import androidx.room.ColumnInfo
import androidx.room.Entity
import androidx.room.PrimaryKey

@Entity(tableName = "suggestionCacheItems")
data class SuggestionCacheItem (
    @PrimaryKey(autoGenerate = true) val id: Int,

    @ColumnInfo(name = "type")
    val type: String,

    @ColumnInfo(name = "text")
    val text: String,

    @ColumnInfo(name = "created")
    val created: Long = System.currentTimeMillis()
)
