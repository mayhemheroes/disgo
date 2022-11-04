package fuzz

import "strconv"
import "github.com/disgoorg/disgo/discord"

func mayhemit(bytes []byte) int {

    var num int
    if len(bytes) > 1 {
        num, _ = strconv.Atoi(string(bytes[0]))

        switch num {
    
        case 0:
            content := string(bytes)

            length := len(content)
            str1 := content[0:length/2]
            str2 := content[length/2:length-1]
            discord.UserTag(str1, str2)
            return 0

        case 1:
            content := string(bytes)

            length := len(content)
            str1 := content[0:length/2]
            str2 := content[length/2:length-1]

            var test discord.EmbedBuilder
            test.AddField(str1, str2, true)
            return 0

        case 2:
            content := string(bytes)

            var test discord.EmbedBuilder
            test.SetAuthorURLf(content)
            return 0

        case 3:
            content := string(bytes)

            var test discord.EmbedBuilder
            test.SetImage(content)
            return 0

        case 4:
            content := string(bytes)

            var test discord.EmbedBuilder
            test.SetTitle(content)
            return 0

        default:
            content := string(bytes)
            discord.InviteURL(content)
            return 0

        }
    }
    return 0
}

func Fuzz(data []byte) int {
    _ = mayhemit(data)
    return 0
}