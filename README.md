Sideproject for testing prupose

This repository provides a Docker Hub pull-through cache that helps speed up Docker image pulls by caching images locally. The current implementation supports local storage, where images pulled from Docker Hub are cached and can be reused by other Docker clients. 


## Features
- Caches Docker images pulled from Docker Hub.
- Supports local storage of cached images.
- Improves the speed of subsequent pulls of the same images and avoid getting docker hub rate-limits

### Future Enhancements
- **Google Container Registry (GCR) Support**: Future support for Google Container Registry (GCR) will allow for caching of images from GCR.
- **Amazon S3 Storage**: We plan to add support for storing cached images in Amazon S3 to ensure a more scalable and durable caching solution across multiple instances.

