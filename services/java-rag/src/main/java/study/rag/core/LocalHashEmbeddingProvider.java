package study.rag.core;

import org.springframework.stereotype.Component;

import java.nio.ByteBuffer;
import java.nio.charset.StandardCharsets;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;
import java.util.Locale;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

@Component
public class LocalHashEmbeddingProvider {
  private static final Pattern TOKEN_PATTERN = Pattern.compile("[a-zA-Z0-9_]+");
  private static final int DIMENSIONS = 128;

  public double[] embed(String text) {
    double[] vector = new double[DIMENSIONS];
    Matcher matcher = TOKEN_PATTERN.matcher(text.toLowerCase(Locale.ROOT));
    while (matcher.find()) {
      byte[] digest = sha256(matcher.group());
      int index = Math.floorMod(ByteBuffer.wrap(digest, 0, 4).getInt(), DIMENSIONS);
      double sign = digest[4] % 2 == 0 ? 1.0 : -1.0;
      vector[index] += sign;
    }

    double length = 0.0;
    for (double value : vector) {
      length += value * value;
    }
    length = Math.sqrt(length);
    if (length == 0.0) {
      return vector;
    }
    for (int i = 0; i < vector.length; i++) {
      vector[i] = vector[i] / length;
    }
    return vector;
  }

  private byte[] sha256(String token) {
    try {
      return MessageDigest.getInstance("SHA-256").digest(token.getBytes(StandardCharsets.UTF_8));
    } catch (NoSuchAlgorithmException exception) {
      throw new IllegalStateException("SHA-256 is not available", exception);
    }
  }
}
