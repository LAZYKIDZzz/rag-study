package study.rag.core;

import org.springframework.stereotype.Repository;
import study.rag.core.model.Chunk;
import study.rag.core.model.Document;

import java.time.Instant;
import java.util.ArrayList;
import java.util.Comparator;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

@Repository
public class DocumentRepository {
  private final Map<String, Document> documents = new HashMap<>();
  private final Map<String, List<Chunk>> chunksByDocument = new HashMap<>();

  public synchronized Document create(String title, String content, Map<String, Object> metadata) {
    Document document = new Document(Ids.next("doc"), title, content, metadata, false, Instant.now());
    documents.put(document.id(), document);
    return document;
  }

  public synchronized Document get(String documentId) {
    Document document = documents.get(documentId);
    if (document == null) {
      throw new NotFoundException("Document " + documentId + " was not found");
    }
    return document;
  }

  public synchronized List<Document> list() {
    return documents.values().stream()
        .sorted(Comparator.comparing(Document::createdAt).reversed())
        .toList();
  }

  public synchronized void replaceChunks(String documentId, List<Chunk> chunks) {
    Document document = get(documentId);
    documents.put(documentId, document.markIndexed());
    chunksByDocument.put(documentId, new ArrayList<>(chunks));
  }

  public synchronized int chunkCount(String documentId) {
    return chunksByDocument.getOrDefault(documentId, List.of()).size();
  }
}
