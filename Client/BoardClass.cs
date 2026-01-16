using System.Collections.Generic;

public class BoardClass
{
    public List<String> User_Ashtray { get; set; }
    public List<String> Opponent_Ashtray { get; set; }
    public List<String> User_Hand { get; set; }
    public List<String> User_Discard { get; set; }
    public String[] Sealed_Cards = new String[3];
    public String[] User_ConscriptArea = new String[5];
    public String[] Opponent_ConscriptArea = new String[5];
    public String[] User_SpellArea = new String[5];
    public String[] Opponent_SpellArea = new String[5];
    public int User_Health = 1000;
    public int Opponent_Health = 1000;
    public int User_DeckCount = 40;
    public int Opponent_DeckCount = 40;
}