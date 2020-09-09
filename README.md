# Artisan

Artisan is a build orchestrator for mono repositories powered by Docker.
Artisan helps you build complex codebases without the need to install all build tools.
 
## How it works
Let's have a repository with an application written in Java and build by Gradle and some 
Vue frontend. Usually, we need to install the correct version of NPM, Java, and Gradle.

For Artisan, the Java backend and Vue frontend are separated modules with dependency. 
Both modules have a 'build' task. The Artisan executes tasks within an own docker container. 

## Installation
If you're Linux or Mac OS user, you can use the following command:

```bash
curl -sfL https://artisan.unravela.io/install.sh | sh
```
