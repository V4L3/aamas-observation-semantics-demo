package ch.unisg.ics.interactions;

import cartago.OPERATION;
import cartago.ObsProperty;
import jason.asSyntax.ASSyntax;
import org.hyperagents.yggdrasil.cartago.artifacts.HypermediaTDArtifact;
import org.hyperagents.yggdrasil.cartago.artifacts.SyncedHypermediaTDArtifact;

public class SyncedOccupancySensor extends SyncedHypermediaTDArtifact {

  public void init() {
    super.init();
    defineObsProperty("fallDetected", false);
    ObsProperty fallDetected = getObsProperty("fallDetected");
    String artId = getArtifactId().toString();
    fallDetected.addAnnot(ASSyntax.createStructure("roomName", ASSyntax.createNumber(Double.parseDouble(artId.substring(artId.length() - 1)))));

  }

  public void init(Object initializationParameters) {
    super.init(initializationParameters);
    defineObsProperty("fallDetected", false);
    ObsProperty fallDetected = getObsProperty("fallDetected");
    String artId = getArtifactId().toString();
    fallDetected.addAnnot(ASSyntax.createStructure("roomName", ASSyntax.createNumber(Double.parseDouble(artId.substring(artId.length() - 1)))));
  }

  @OPERATION
  public void triggerFallDetected() {
    ObsProperty fallDetected = getObsProperty("fallDetected");
    updateInternalTimestamp();
    fallDetected.updateValue(true);
  }

  @Override
  protected void registerInteractionAffordances() {
    // Register one action affordance with an input schema
    super.registerInteractionAffordances();
    registerActionAffordance("http://example.org/triggerFallDetected", "triggerFallDetected", "triggerFallDetected");
  }

}
