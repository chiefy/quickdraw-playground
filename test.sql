SELECT COUNT(*) AS c FROM (
SELECT Draws.ID, DATE(Draws.drawDate) AS dd 
FROM Draws 
INNER JOIN Picks ON Draws.ID = Picks.drawID 
WHERE Picks.pickNum = 1 AND dd BETWEEN '2017-01-01'  AND '2017-01-02'
GROUP BY Draws.ID
) AS found