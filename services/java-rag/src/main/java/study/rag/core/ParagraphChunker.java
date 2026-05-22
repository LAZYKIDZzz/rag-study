package study.rag.core;

import org.springframework.stereotype.Component;
import study.rag.core.model.Chunk;
import study.rag.core.model.Document;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;

@Component
public class ParagraphChunker {
  private static final int MAX_WORDS = 120;
  private static final int OVERLAP_WORDS = 20;

  public List<Chunk> chunk(Document document, String cleanText) {
    String[] words = cleanText.split("\\s+");
    List<Chunk> chunks = new ArrayList<>();
    if (words.length == 0 || words[0].isBlank()) {
      return chunks;
    }

    int start = 0;
    int ordinal = 0;
    while (start < words.length) {
      int end = Math.min(start + MAX_WORDS, words.length);
      String content = String.join(" ", List.of(words).subList(start, end));
      HashMap<String, Object> metadata = new HashMap<>(document.metadata());
      metadata.put("ordinal", ordinal);
      metadata.put("word_start", start);
      metadata.put("word_end", end);
      chunks.add(new Chunk(Ids.next("chk"), document.id(), document.title(), ordinal, content, metadata));
      if (end == words.length) {
        break;
      }
      start = end - OVERLAP_WORDS;
      ordinal += 1;
    }
    return chunks;
  }
}
