# EvilSukusBot3000

A little Discord that checks whenether [sukus21](https://github.com/sukus21) has
published his programming language

## Backstory

When I was at the Coding Pirates Game Jam 2024 me and sukus talked a lot about
gamedev, osdev and just programming in general. He also showed me his
programming language which is GML-like but it compiles to ARM so you can run it
on a Nintendo DS. Fast forward to the 2025 Game Jam he **still haven't published
the GitHub repo**. So now I have decided to build a Discord bot to check if he
has published it.

## Design

- It should be fast and run in a tiny Docker Container.

  We will optimize the size by using multiple build steps in the Dockerfile.

- It should check his repositories from the GitHub API every 5 minutes.

  We can do this because the GitHub API accepts 60 requests from a single IP
  each hour without authentication.

- It should use the Go Programming Language for speed and tiny executable size.

  I thought for a bit over the language. My initial idea was using Typescript
  and The Deno Runtime but I just decided that I wan't to use Go because it is a
  lot faster than Typescript.

- It should use caching by hashing the GitHub API response.

  This is yet another optimization. We will cache the list of repos and then
  with each request we will compare the hash. If it is the same just skip, but
  if it is diffrent check each repos name for the regexp (The Go Package)
  `(?i)spin`. If it succedes ping on Discord.
