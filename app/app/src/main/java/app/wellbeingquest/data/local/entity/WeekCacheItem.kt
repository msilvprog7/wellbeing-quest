package app.wellbeingquest.data.local.entity

import androidx.room.ColumnInfo
import androidx.room.Entity
import androidx.room.PrimaryKey

@Entity(
    tableName = "weekCacheItems",
    indices = [androidx.room.Index(value = ["name"])])
data class WeekCacheItem (
    @PrimaryKey(autoGenerate = true) val id: Int = 0,

    @ColumnInfo(name = "name")
    val name: String,

    @ColumnInfo(name = "feeling")
    val feeling: String,

    @ColumnInfo(name = "activity")
    val activity: String,

    @ColumnInfo(name = "created")
    val created: Long = System.currentTimeMillis()
)
