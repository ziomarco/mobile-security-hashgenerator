import 'dart:io';

import 'package:encrypt/encrypt.dart';

void main() async {
  final file = File(
          '/Users/ziomarco/Documents/Projects/Personal/mobile-security-hashgenerator/encrypted.json')
      .readAsBytesSync();

  var keyString = 'k4rAN45oL8LxH21wX2nRTDB5o1uYnnrB';

  final key = Key.fromUtf8(keyString);

  final iv = IV(file.sublist(0, 16));

  final encrypter = Encrypter(AES(key, mode: AESMode.cbc, padding: 'PKCS7'));

  final encrypted = Encrypted(file);

  final decrypted = encrypter.decrypt(encrypted, iv: iv);

  final decryptedFileAsString =
      String.fromCharCodes(decrypted.codeUnits.sublist(16));
  print(decryptedFileAsString);
}
