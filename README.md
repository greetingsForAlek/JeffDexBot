# The JeffDex Bot

The official discord bot for JeffDex. JeffDex is a catalogue of characters from Jeff.

This guide will show you how to use this bot.

## General Use
The bot comes with many built in slash commands. The whole point of the bot is a fun guessing game where users have to guess the character that is shown. But it also includes a daily shop where users can buy characters as well as hourly drops where a character can appear every hour.

### The Ping Command
`/ping`
The most basic command in the whole bot. It's mainly used for testing the bot and all it does is have the bot reply with "pong".

### The Guess Command
`/guess`
This command initiates the guessing game. The first user to guess the correct name of the character is awarded the character, and the character is added to their index (accessed via the [index](#the-index-command) command)

When the command is run it messages an embedded message with an image of the character and also the character's [rarity](#rarities--shinies). It also has a button that is labelled `Who is this?` and when pressed the button will reveal the name of the character, making the game easier for user that may not be familiar with the characters.

Characters also have a chance of being shiny. Shiny characters award a coin reward depending on the rarity of the character. [See the rarities and shinies section](#rarities--shinies) For info on shinies.

***NOTE:*** If a user already owns a character, they will be awarded coins for a duplicate. Please see the [rewards](#rewards) section for more info on duplicate and shiny rewards.

### The Index Command
`/index` When this command is run, users can view their character index, showing all of their collected characters listed by [rarity](#rarities--shinies).

Any character not collected is represented by "❓❓❓" and the character's name is not revealled, although the rarity of the character is revealled.

***Note:*** The bot responds to the index command with an *ephemeral* message, therefore other user cannot see your character index.

### The Wallet Command
`/wallet` This command allows users to check their coin balance. The bot responds with a simple ephemeral message so other users cannot see.

Coins are collected from shinies and duplicates. See the [rewards section](#rewards) for more info.

### The Shop Command
`/shop` When this command is run, the bot responds with an embedded message showcasing all characters currently in rotation for the daily shop.

Characters in rotation reset every 24 hours, and only **Rare, Epic and Legendary** characters can appear in the shop. Prices depend on the rarity of the character. (For more info on rarities please see the [rarities and shinies section](#rarities--shinies)).

The embed also includes buttons to purchase each of the characters. Once the button is clicked, if you can afford to buy the character, the character is added to your index and the price in coins is subtracted.

***NOTE:*** If you already own the character that is in the shop, **No duplicate reward will be given**, so please check your index before purchasing characters, otherwise you'll waste your hard-earned cash ;)

## Rarities + Shinies
All of the rarities currently used by the bot and their chances of appearing are:

- Common (50% chance of appearing)
- Uncommon (25% chance of appearing)
- Rare (15% chance of appearing)
- Epic (7% chance of appearing)
- Legendary (3% chance of appearing)

Each time a round is initiated, the character has a **15% chance** of being **shiny**. Shinies give [rewards](#rewards) depending on the rarity.

## Rewards
Here are the **duplicate reward amounts** for each rarity:

- Common (10 coins)
- Uncommon (20 coins)
- Rare (50 coins)
- Epic (100 coins)
- Legendary (250 coins)

And here are the **shiny reward amounts** for each rarity:

- Common (25 coins)
- Uncommon (50 coins)
- Rare (100 coins)
- Epic (250 coins)
- Legendary (500 coins)