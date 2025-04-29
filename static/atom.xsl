<xsl:stylesheet version="1.0"
                xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
                xmlns:atom="http://www.w3.org/2005/Atom"
                exclude-result-prefixes="atom">

<xsl:output method="html" encoding="UTF-8" indent="yes"/>

<xsl:template match="/">
  <html>
    <head>
      <title><xsl:value-of select="/atom:feed/atom:title"/></title>
      <style>
        body { font-family: sans-serif; line-height: 1.6; padding: 20px; max-width: 800px; margin: auto; }
        h1, h2 { color: #333; }
        .entry { border-bottom: 1px solid #eee; padding-bottom: 15px; margin-bottom: 15px; }
        .entry h3 { margin-bottom: 5px; }
        .meta { font-size: 0.9em; color: #666; }
        a { color: #007bff; text-decoration: none; }
        a:hover { text-decoration: underline; }
        code { background-color: #f8f9fa; padding: 2px 4px; border-radius: 4px; font-size: 1rem; }
      </style>
    </head>
    <body>
      <h1><xsl:value-of select="/atom:feed/atom:title"/></h1>
      <p><xsl:value-of select="/atom:feed/atom:subtitle"/></p>
      <hr/>
      <h2>¿Cómo suscribirme a este feed?</h2>
        <p>Para suscribirte a este feed, puedes usar un lector de feeds RSS o Atom. Los hay disponibles gratuitamente y de pago, y de escritorio o online. Aquí tienes los que recomendamos:</p>
        <ul>
            <li><a href="https://www.feedly.com/" target="_blank">Feedly</a></li>
            <li><a href="https://www.inoreader.com/" target="_blank">Inoreader</a></li>

            <li><a href="https://thunderbird.net/" target="_blank">Thunderbird</a></li>
            <li><a href="https://www.rssowl.org/" target="_blank">RSSOwl</a></li>
            <li><a href="https://www.rssreader.com/" target="_blank">RSS Reader</a></li>
        </ul>
        <p>Una vez instalado o registrado, basta con copiar y pegar la URL de este feed en el lector de feeds. La URL es:</p>
        <p><code><xsl:value-of select="/atom:feed/atom:link[@rel='self']/@href"/></code></p>

      <hr/>
      <h2>Últimas entradas</h2>
      <xsl:apply-templates select="/atom:feed/atom:entry"/>
    </body>
  </html>
</xsl:template>

<xsl:template match="atom:entry">
  <div class="entry">
    <h3>
      <xsl:variable name="entryLink">
         <xsl:choose>
            <xsl:when test="atom:link[@rel='alternate']/@href">
               <xsl:value-of select="atom:link[@rel='alternate']/@href"/>
            </xsl:when>
            <xsl:otherwise>
               <xsl:value-of select="atom:link[not(@rel)]/@href"/>
            </xsl:otherwise>
         </xsl:choose>
      </xsl:variable>
      <a href="{$entryLink}" target="_blank">
        <xsl:value-of select="atom:title"/>
      </a>
    </h3>
    <div class="meta">
      Publicado el: <xsl:value-of select="substring(atom:published, 1, 10)"/>
      <xsl:if test="atom:author/atom:name">
        | Autor: <xsl:value-of select="atom:author/atom:name"/>
      </xsl:if>
    </div>
    <xsl:if test="atom:summary">
       <p><xsl:value-of select="atom:summary" disable-output-escaping="yes"/></p>
    </xsl:if>
  </div>
</xsl:template>

</xsl:stylesheet>