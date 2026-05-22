package study.rag.core;

import org.springframework.stereotype.Repository;
import study.rag.core.model.Chunk;
import study.rag.core.model.SearchResult;

import java.util.Comparator;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

@Repository
public class InMemoryVectorStore {
  private final Map<String, VectorRecord> records = new HashMap<>();

  public synchronized void deleteByDocument(String documentId) {
    records.entrySet().removeIf(entry -> entry.getValue().chunk().documentId().equals(documentId));
  }

  public synchronized void upsertMany(List<Chunk> chunks, List<double[]> embeddings) {
    for (int i = 0; i < chunks.size(); i++) {
      records.put(chunks.get(i).id(), new VectorRecord(chunks.get(i), embeddings.get(i)));
    }
  }

  public synchronized List<SearchResult> search(double[] queryEmbedding, int topK, Map<String, Object> filters) {
    return records.values().stream()
        .filter(record -> matches(record.chunk(), filters))
        .map(record -> new SearchResult(record.chunk(), cosine(queryEmbedding, record.embedding())))
        .sorted(Comparator.comparingDouble(SearchResult::score).reversed())
        .limit(topK)
        .toList();
  }

  public synchronized int countForDocument(String documentId) {
    return (int) records.values().stream()
        .filter(record -> record.chunk().documentId().equals(documentId))
        .count();
  }

  private boolean matches(Chunk chunk, Map<String, Object> filters) {
    if (filters == null || filters.isEmpty()) {
      return true;
    }
    return filters.entrySet().stream()
        .allMatch(entry -> entry.getValue().equals(chunk.metadata().get(entry.getKey())));
  }

  private double cosine(double[] left, double[] right) {
    double dot = 0.0;
    double leftLength = 0.0;
    double rightLength = 0.0;
    for (int i = 0; i < left.length; i++) {
      dot += left[i] * right[i];
      leftLength += left[i] * left[i];
      rightLength += right[i] * right[i];
    }
    if (leftLength == 0.0 || rightLength == 0.0) {
      return 0.0;
    }
    return dot / (Math.sqrt(leftLength) * Math.sqrt(rightLength));
  }

  private record VectorRecord(Chunk chunk, double[] embedding) {
  }
}
