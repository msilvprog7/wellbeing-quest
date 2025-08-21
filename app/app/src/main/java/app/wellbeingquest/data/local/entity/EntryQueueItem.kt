package app.wellbeingquest.data.local.entity

import androidx.room.ColumnInfo
import androidx.room.Entity
import androidx.room.PrimaryKey

@Entity(tableName = "entryQueueItems")
data class EntryQueueItem (
    @PrimaryKey(autoGenerate = true) val id: Int,

    @ColumnInfo(name = "activity")
    val activity: String,

    @ColumnInfo(name = "feelings")
    val feelings: List<String>,

    @ColumnInfo(name = "created")
    val created: Long = System.currentTimeMillis()
)