package study.rag.core.dto;

import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.List;

public record RetrievalResponse(
    String query,
    @JsonProperty("rewritten_query") String rewrittenQuery,
    @JsonProperty("retrieved_chunks") List<RetrievedChunk> retrievedChunks) {
}
