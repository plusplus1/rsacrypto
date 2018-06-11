#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""
Date: 2018/04/25 13:22:18
"""

import base64

import Crypto.Hash.SHA
import Crypto.PublicKey.RSA
import Crypto.Signature.PKCS1_v1_5
from Crypto.Cipher import PKCS1_v1_5


def _to_string(val):
    if isinstance(val, str):
        return val
    if isinstance(val, unicode):
        return val.encode('utf8')
    if val is None:
        return ""
    return str(val)


class RSAUtil(object):
    """
        RSA 工具集
    """

    _decrypt_slice_size = 128
    _encrypt_slice_size = 100

    _rsa_obj = getattr(Crypto.PublicKey.RSA, '_RSAobj')

    @classmethod
    def load_key(cls, str_key):
        """
        加载公私钥
        :param str_key:
        :return:
        """
        return Crypto.PublicKey.RSA.importKey(_to_string(str_key))

    @classmethod
    def _pretreat_key(cls, key):
        if isinstance(key, (bytes, str, unicode)):
            return cls.load_key(key)
        return key

    @classmethod
    def pub_encrypt(cls, text, key):
        """
        使用公钥加密数据
        :param text:
        :param key:
        :return:
        """
        key = cls._pretreat_key(key)
        assert isinstance(key, cls._rsa_obj)

        crypt_text = _to_string(text)
        out_text = ""

        pkcs_obj = PKCS1_v1_5.new(key)
        while crypt_text:
            tmp_text = crypt_text[:cls._encrypt_slice_size]
            crypt_text = crypt_text[cls._encrypt_slice_size:]
            out_text += pkcs_obj.encrypt(tmp_text)
            pass

        return base64.b64encode(out_text)

    @classmethod
    def pri_decrypt(cls, text, key):
        """
        使用私钥解密数据
        :param text:
        :param key:
        :return:
        """
        key = cls._pretreat_key(key)
        assert isinstance(key, cls._rsa_obj)

        crypt_text = base64.b64decode(_to_string(text))
        out_text = ""

        pcks_obj = PKCS1_v1_5.new(key)

        while crypt_text:
            tmp_text = crypt_text[:cls._decrypt_slice_size]
            crypt_text = crypt_text[cls._decrypt_slice_size:]
            out_text += pcks_obj.decrypt(tmp_text, "")

        return out_text

    @classmethod
    def calc_signature(cls, text, pri_key):
        """
        使用私钥进行数字签名
        :param text:
        :param pri_key:
        :return:
        """

        pri_key = cls._pretreat_key(pri_key)
        signer = Crypto.Signature.PKCS1_v1_5.new(pri_key)
        r = signer.sign(Crypto.Hash.SHA.new(_to_string(text)))
        return base64.b64encode(r)

    @classmethod
    def verify_signature(cls, text, sign, pub_key):
        """
        使用公钥校验数字签名
        :param text:
        :param sign:
        :param pub_key:
        :return:
        """
        pub_key = cls._pretreat_key(pub_key)
        signer = Crypto.Signature.PKCS1_v1_5.new(pub_key)
        return signer.verify(
            Crypto.Hash.SHA.new(_to_string(text)),
            base64.b64decode(sign)
        )

    pass
