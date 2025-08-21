package app.wellbeingquest.data.local.database

import androidx.room.Database
import androidx.room.RoomDatabase
import androidx.room.TypeConverters
import app.wellbeingquest.data.local.dao.EntryDraftDao
import app.wellbeingquest.data.local.dao.EntryQueueItemDao
import app.wellbeingquest.data.local.dao.SuggestionCacheItemDao
import app.wellbeingquest.data.local.dao.WeekCacheItemDao
import app.wellbeingquest.data.local.entity.Converters
import app.wellbeingquest.data.local.entity.EntryDraft
import app.wellbeingquest.data.local.entity.EntryQueueItem
import app.wellbeingquest.data.local.entity.SuggestionCacheItem
import app.wellbeingquest.data.local.entity.WeekCacheItem

@Database(
    entities = [EntryDraft::class, EntryQueueItem::class, SuggestionCacheItem::class, WeekCacheItem::class],
    version = 1
)
@TypeConverters(Converters::class)
abstract class AppDatabase : RoomDatabase() {
    abstract fun entryDraftDao(): EntryDraftDao
    abstract fun entryQueueItemDao(): EntryQueueItemDao
    abstract fun suggestionCacheItemDao(): SuggestionCacheItemDao
    abstract fun weekCacheItemDao(): WeekCacheItemDao
}
