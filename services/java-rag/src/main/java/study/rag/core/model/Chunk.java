package study.rag.core.model;

import java.util.Map;

public record Chunk(
    String id,
    String documentId,
    String title,
    int ordinal,
    String content,
    Map<String, Object> metadata) {
}
