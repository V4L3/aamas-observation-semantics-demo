# Use a base image with Java 21 installed
FROM eclipse-temurin:21-jdk

# Set the working directory inside the container
WORKDIR /app

# Copy the entire project into the container
COPY . .

# Set execute permissions for gradlew
RUN chmod +x ./gradlew

# Command to run the application
CMD ["./gradlew", "run", "--no-daemon"]

# Expose port 8080
EXPOSE 8080
