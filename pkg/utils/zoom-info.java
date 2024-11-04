import java.io.OutputStream;
import java.net.HttpURLConnection;
import java.net.URL;
import java.nio.charset.StandardCharsets;

public class ShareData {

    public static void main(String[] args) {
        try {
            // Set the API endpoint URL for ZoomInfo or your target service
            URL url = new URL("https://api.zoominfo.com/v2/endpoint"); // Replace with actual ZoomInfo endpoint
            HttpURLConnection connection = (HttpURLConnection) url.openConnection();

            // Set request type and headers
            connection.setRequestMethod("POST");
            connection.setRequestProperty("Content-Type", "application/json");
            connection.setRequestProperty("Authorization", "Bearer YOUR_API_TOKEN"); // Replace with your API token
            connection.setDoOutput(true);

            // Prepare JSON payload with IP information
            String jsonInputString = "{ \"ip\": \"123.45.67.89\" }"; // Replace with your IP info

            // Send the JSON data
            try (OutputStream os = connection.getOutputStream()) {
                byte[] input = jsonInputString.getBytes(StandardCharsets.UTF_8);
                os.write(input, 0, input.length);
            }

            // Check for success response
            int responseCode = connection.getResponseCode();
            if (responseCode == HttpURLConnection.HTTP_OK) {
                System.out.println("Data sent successfully to ZoomInfo!");
            } else {
                System.out.println("Failed to send data. Response Code: " + responseCode);
            }

            connection.disconnect();

        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}
