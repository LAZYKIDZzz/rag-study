package study.rag.core;

import java.util.UUID;

final class Ids {
  private Ids() {
  }

  static String next(String prefix) {
    return prefix + "_" + UUID.randomUUID().toString().replace("-", "").substring(0, 12);
  }
}
