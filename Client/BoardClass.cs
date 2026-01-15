using System.Collections.Generic;

public class BoardClass
{
    public List<CardClass> User_Ashtray { get; set; }
    public List<CardClass> Opponent_Ashtray { get; set; }
    public List<CardClass> User_Hand { get; set; }
    public List<CardClass> User_Discard { get; set; }
    public CardClass[] Sealed_Cards = new CardClass[3];
    public CardClass[] User_ConscriptArea = new CardClass[5];
    public CardClass[] Opponent_ConscriptArea = new CardClass[5];
    public CardClass[] User_SpellArea = new CardClass[5];
    public CardClass[] Opponent_SpellArea = new CardClass[5];
    public int User_Health = 1000;
    public int Opponent_Health = 1000;
    public int User_DeckCount = 40;
    public int Opponent_DeckCount = 40;
}