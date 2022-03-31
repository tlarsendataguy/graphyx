import 'dart:io';

import 'package:flutter_test/flutter_test.dart';

Future generateFontFile(String source, String target, String varName) async {
  var uri = Uri(path: source);
  var file = File.fromUri(uri);
  var bytes = file.readAsBytesSync();
  var printFrom = 0;

  var log = File.fromUri(Uri(path: target));
  var sink = log.openWrite();
  sink.writeln("List<int> $varName = [");
  while (printFrom < bytes.length) {
    var printTo = printFrom + 50;
    if (printTo > bytes.length) {
      printTo = bytes.length;
    }
    var bytesStr = bytes.sublist(printFrom, printTo).toString();
    bytesStr = bytesStr.replaceAll(RegExp(r'[\[\]]'), "");
    sink.writeln("$bytesStr,");
    printFrom = printTo;
  }
  sink.writeln("];");
  await sink.flush();
  await sink.close();
}

main(){
  test('generate mono_font_file.dart', () async {
    await generateFontFile(
      """fonts/JetBrainsMono-Regular.ttf""",
      """lib/mono_font_file.dart""",
      "monoFontFile",
    );
  });

  test('generate material_icons.dart', () async {
      await generateFontFile(
        """fonts/MaterialIcons-Regular.ttf""",
        """lib/material_icons.dart""",
        "materialIcons",
      );
  });
}
