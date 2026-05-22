package study.rag.core.model;

import java.time.Instant;
import java.util.Map;

public record Document(
    String id,
    String title,
    String content,
    Map<String, Object> metadata,
    boolean indexed,
    Instant createdAt) {
  public Document markIndexed() {
    return new Document(id, title, content, metadata, true, createdAt);
  }
}
