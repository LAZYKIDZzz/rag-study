package study.rag.core.model;

import java.time.Instant;
import java.util.List;

public record ChatMessage(
    String role,
    String content,
    Instant createdAt,
    List<SearchResult> citations) {
}
