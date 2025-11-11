package app.wellbeingquest.data.service.api

import android.security.keystore.KeyGenParameterSpec
import android.security.keystore.KeyProperties
import java.security.KeyStore
import javax.crypto.Cipher
import javax.crypto.KeyGenerator
import javax.crypto.SecretKey
import javax.crypto.spec.GCMParameterSpec
import android.util.Base64
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class CryptoManager @Inject constructor() {

    companion object {
        private const val ANDROID_KEYSTORE = "AndroidKeyStore"
        private const val KEY_ALIAS = "session_token_key"
        private const val AES_MODE = "AES/GCM/NoPadding"
        private const val IV_SIZE = 12
        private const val TAG_SIZE = 128
    }

    private val keyStore = KeyStore.getInstance(ANDROID_KEYSTORE).apply { load(null) }

    private fun getSecretKey(): SecretKey {
        if (!keyStore.containsAlias(KEY_ALIAS)) {
            val keyGen = KeyGenerator.getInstance(KeyProperties.KEY_ALGORITHM_AES, ANDROID_KEYSTORE)
            val spec = KeyGenParameterSpec.Builder(KEY_ALIAS,
                KeyProperties.PURPOSE_ENCRYPT or KeyProperties.PURPOSE_DECRYPT)
                .setBlockModes(KeyProperties.BLOCK_MODE_GCM)
                .setEncryptionPaddings(KeyProperties.ENCRYPTION_PADDING_NONE)
                .build()
            keyGen.init(spec)
            keyGen.generateKey()
        }
        return (keyStore.getEntry(KEY_ALIAS, null) as KeyStore.SecretKeyEntry).secretKey
    }

    fun encrypt(plainText: String): String {
        val cipher = Cipher.getInstance(AES_MODE)
        cipher.init(Cipher.ENCRYPT_MODE, getSecretKey())
        val iv = cipher.iv
        val encrypted = cipher.doFinal(plainText.toByteArray())
        val combined = iv + encrypted
        return Base64.encodeToString(combined, Base64.DEFAULT)
    }

    fun decrypt(cipherText: String): String {
        val combined = Base64.decode(cipherText, Base64.DEFAULT)
        val iv = combined.sliceArray(0 until IV_SIZE)
        val encrypted = combined.sliceArray(IV_SIZE until combined.size)
        val cipher = Cipher.getInstance(AES_MODE)
        val spec = GCMParameterSpec(TAG_SIZE, iv)
        cipher.init(Cipher.DECRYPT_MODE, getSecretKey(), spec)
        return String(cipher.doFinal(encrypted))
    }
}